import { async, TestBed } from '@angular/core/testing';
import { AuthGuard } from './auth-guard.service';
import { RouterTestingModule } from '@angular/router/testing';
import {ToastrModule} from 'ngx-toastr';

describe('AuthGuardService', () => {
  let service: AuthGuard;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        RouterTestingModule.withRoutes([]),
        ToastrModule.forRoot(),
      ],
      declarations: []
    });
    service = TestBed.inject(AuthGuard);
  }));

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
