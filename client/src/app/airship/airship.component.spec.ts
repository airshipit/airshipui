import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AirshipComponent } from './airship.component';

describe('AirshipComponent', () => {
  let component: AirshipComponent;
  let fixture: ComponentFixture<AirshipComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AirshipComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AirshipComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
