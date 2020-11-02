/*
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

import { Injectable } from '@angular/core';
import { Router, CanActivate, Event as RouterEvent, NavigationStart, NavigationEnd, NavigationCancel, NavigationError } from '@angular/router';
import { Log } from 'src/services/log/log.service';
import { LogMessage } from 'src/services/log/log-message';
import { WsService } from 'src/services/ws/ws.service';
import { WsReceiver, WsMessage, WsConstants } from 'src/services/ws/ws.models';

@Injectable({
  providedIn: 'root'
})

export class AuthGuard implements WsReceiver, CanActivate {
  // static router for those who may need it, I'm looking at your app components
  public static router: Router;

  private className = this.constructor.name;
  private loading = false;
  private sendToLogin = false;

  type = WsConstants.UI;
  component = WsConstants.AUTH;

  // Called by the logout link at the top right of the page
  public static logout(): void {
    // blank out the object storage so we can't get re authenticate
    WsService.token = undefined;
    WsService.refreshToken = undefined;

    // blank out the local storage so we can't get re authenticate
    localStorage.removeItem('airshipUI-token');

    // turn off the log panel, no logs for you!
    AuthGuard.toggleLogPanel(false);

    // best to begin at the beginning so send the user back to /login
    this.router.navigate(['/login']);
  }

  // flip the log panel according to where we are in the world
  public static toggleLogPanel(authenticated): void {
    const accordion = document.getElementById('logAccordion');
    if (authenticated && accordion.style.display === 'none') {
      accordion.style.display = '';
    } else if (!authenticated) {
      accordion.style.display = 'none';
    }
  }

  constructor(private websocketService: WsService, private router: Router) {
    // create a static router so other components can access it if needs be
    AuthGuard.router = router;

    this.websocketService.registerFunctions(this);
    // listen to the evens that are sent out from the angular router so we don't wind up in an endless loop
    this.router.events.subscribe((e: RouterEvent) => {
      this.navigationInterceptor(e);
    });
  }

  async receiver(message: WsMessage): Promise<void> {
    if (message.hasOwnProperty(WsConstants.ERROR)) {
      Log.Error(new LogMessage('Error received in AuthGuard', this.className, message));
      AuthGuard.logout();
    } else {
      switch (message.subComponent) {
        case WsConstants.APPROVED:
          this.setToken(message.token, false);
          Log.Debug(new LogMessage('Auth approved received', this.className, message));
          // redirect to / only when on /login otherwise leave the path where it was before the auth attempt
          const location = window.location.pathname;
          if (location === '/login' || location === '/login/') {
            this.router.navigate(['/']);
          }
          break;
        case WsConstants.DENIED:
          AuthGuard.logout();
          Log.Debug(new LogMessage('Auth denied received', this.className, message));
          break;
        case WsConstants.REFRESH:
          this.setToken(message.refreshToken, true);
          Log.Debug(new LogMessage('Auth token refresh received', this.className, message));
          break;
        default:
          AuthGuard.logout();
          Log.Debug(new LogMessage('Unknown auth message received', this.className, message));
          break;
      }
    }
  }

  // this decides if you can show a page
  // TODO: maybe RBAC type of stuff may need to go here
  canActivate(): boolean {
    const authenticated = this.validateToken();
    const location = window.location.pathname;

    // redirect everything to /login if not authenticated
    if (!authenticated && location !== '/login/') {
      // TODO: store the reference url and redirect after login
      // let the loading function complete before sending to login otherwise the redirect fails
      if (this.loading) {
        this.sendToLogin = true;
      } else {
        // loading is complete just send to login
        this.router.navigate(['/login']);
      }
      return true;
    }

    // login page specific details
    // redirect /login to / if authenticated and landing on /login
    // If this is not here when you refresh on the login page it somehow redirects to /
    if (location === '/login/') {
      if (authenticated) {
        this.router.navigate(['/']);
        return false;
      } else {
        return true;
      }
    }

    // flip the link if we're in or out of the fold
    this.toggleAuthButton(authenticated);

    // flip the visibility of the log panel depending on the disposition of the user
    AuthGuard.toggleLogPanel(authenticated);

    return authenticated;
  }

  // flip the text of the login / logout button according to where we are in the world
  private toggleAuthButton(authenticated): void {
    const button = document.getElementById('loginButton');
    const text = button.innerText;
    if (authenticated && text === 'Login') {
      button.innerText = 'Logout';
    } else if (!authenticated && text === 'Logout') {
      button.innerText = 'Login';
    }
  }

  // retrieve the stored token & send it to the go backend for validation
  private getStoredToken(): void {
    const tokenString = localStorage.getItem('airshipUI-token');
    const token = JSON.parse(tokenString);
    if (token !== null) {
      if (token.hasOwnProperty('token')) {
        WsService.token = token.token;
      }
    }
  }

  // the UI frontend is not the decider, the back end is.  If this token is good we continue, if it's not we stop
  private validateToken(): boolean {
    if (WsService.token === undefined) { this.getStoredToken(); }

    // even after all this it's possible to have nothing.  I started with nothing and still have most of it left
    if (WsService.token !== undefined) {
      const message = new WsMessage(this.type, this.component, WsConstants.VALIDATE);
      message.token = WsService.token;

      // if we have a refresh token we also need to include that in the validity check
      if (WsService.refreshToken !== undefined) {
        message.refreshToken = WsService.refreshToken;
      }

      this.websocketService.sendMessage(message);
    }

    return WsService.token !== undefined;
  }

  // store the token locally so we can be authenticated between runs
  private setToken(token, isRefresh): void {
    // set the token for auth check going forward
    if (isRefresh) {
      WsService.refreshToken = token;
    } else {
      WsService.token = token;

      // set the token locally to have a login till browser exits
      const json: any = { token: WsService.token };
      localStorage.setItem('airshipUI-token', JSON.stringify(json));
    }
  }

  // detect navigation events in case we redirect from authguard which would happen too fast to protect /login and cause an endless loop
  // Random Shack Data Processing Dictionary: Endless Loop: n., see Loop, Endless. Loop, Endless: n., see Endless Loop
  private navigationInterceptor(event: RouterEvent): void {
    if (event instanceof NavigationStart) {
      this.loading = true;
    }
    if (event instanceof NavigationEnd) {
      this.loading = false;
      if (this.sendToLogin) {
        this.router.navigate(['/login']);
        this.sendToLogin = false;
      }
    }
    if (event instanceof NavigationCancel) {
      this.loading = false;
    }
    if (event instanceof NavigationError) {
      this.loading = false;
    }
  }
}
