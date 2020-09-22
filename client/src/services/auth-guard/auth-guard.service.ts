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
import { WebsocketService } from 'src/services/websocket/websocket.service';
import { WSReceiver, WebsocketMessage } from 'src/services/websocket/websocket.models';

@Injectable({
  providedIn: 'root'
})

export class AuthGuard implements WSReceiver, CanActivate {
  // static router for those who may need it, I'm looking at your app components
  public static router: Router;

  private className = this.constructor.name;
  private loading = false;
  private sendToLogin = false;
  type = 'ui';
  component = 'auth';

  // Called by the logout link at the top right of the page
  public static logout(): void {
    // blank out the object storage so we can't get re authenticate
    WebsocketService.token = undefined;
    WebsocketService.tokenExpiration = 0;

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

  constructor(private websocketService: WebsocketService, private router: Router) {
    // create a static router so other components can access it if needs be
    AuthGuard.router = router;

    this.websocketService.registerFunctions(this);
    // listen to the evens that are sent out from the angular router so we don't wind up in an endless loop
    this.router.events.subscribe((e: RouterEvent) => {
      this.navigationInterceptor(e);
    });
  }

  async receiver(message: WebsocketMessage): Promise<void> {
    if (message.hasOwnProperty('error')) {
      Log.Error(new LogMessage('Error received in AuthGuard', this.className, message));
      this.websocketService.printIfToast(message);
      AuthGuard.logout();
    } else {
      switch (message.subComponent) {
        case 'approved':
          Log.Debug(new LogMessage('Auth approved received', this.className, message));
          this.setToken(message.token);
          this.router.navigate(['/']);
          break;
        case 'denied':
          Log.Debug(new LogMessage('Auth denied received', this.className, message));
          AuthGuard.logout();
          break;
        default:
          Log.Debug(new LogMessage('Unknown auth message received', this.className, message));
          AuthGuard.logout();
          break;
      }
    }
  }

  // this decides if you can show a page
  // TODO: maybe RBAC type of stuff may need to go here
  canActivate(): boolean {
    const location = window.location.pathname;
    const authenticated = this.isAuthenticated();

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
    // TODO (aschiefe): not super happy about this setup, may need to simplify
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



  // test the auth token to see if we can let the user see the page
  // TODO: maybe RBAC type of stuff may need to go here
  private isAuthenticated(): boolean {
    if (WebsocketService.token === undefined) { this.getStoredToken(); }
    try {
      let authenticated = false;
      // test for token expiration
      // if the token is null the date test will always return true
      if (WebsocketService.token !== undefined && WebsocketService.tokenExpiration > 0) {
        authenticated = WebsocketService.tokenExpiration >= new Date().getTime();
      }
      return authenticated;
    } catch (ex) {
      return false;
    }
  }

  // retrieve the stored token & send it to the go backend for validation
  private getStoredToken(): void {
    const tokenString = localStorage.getItem('airshipUI-token');
    const token = JSON.parse(tokenString);
    if (token !== null) {
      if (token.hasOwnProperty('token')) {
        WebsocketService.token = token.token;
      }
      if (token.hasOwnProperty('date')) {
        WebsocketService.tokenExpiration = token.date;
      }

      // even after all this it's possible to have nothing.  I started with nothing and still have most of it left
      if (WebsocketService.token !== undefined) {
        this.validateToken();
      }
    }
  }

  // the UI frontend is not the decider, the back end is.  If this token is good we continue, if it's not we stop
  private validateToken(): void {
    const message = new WebsocketMessage(this.type, this.component, 'validate');
    message.token = WebsocketService.token;
    this.websocketService.sendMessage(message);
  }

  // store the token locally so we can be authenticated between runs
  private setToken(token): void {
    // calculate 1 hour expiration
    const date = new Date();
    date.setTime(date.getTime() + (1 * 60 * 60 * 1000));

    // set the token for auth check going forward
    WebsocketService.token = token;
    WebsocketService.tokenExpiration = date.getTime();

    // set the token locally to have a login till browser exits
    const json = { date: WebsocketService.tokenExpiration, token: WebsocketService.token };
    localStorage.setItem('airshipUI-token', JSON.stringify(json));
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
