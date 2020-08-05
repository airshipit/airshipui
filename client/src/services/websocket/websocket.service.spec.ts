import {TestBed} from '@angular/core/testing';
import {WebsocketService} from './websocket.service';
import {ToastrModule} from 'ngx-toastr';

describe('WebsocketService', () => {
  let service: WebsocketService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        ToastrModule.forRoot(),
      ]
    });
    service = TestBed.inject(WebsocketService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
