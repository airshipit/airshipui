import { TestBed } from '@angular/core/testing';

import { Log } from './log.service';

describe('LogService', () => {
  let service: Log;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(Log);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
