import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CTLComponent } from './ctl.component';

describe('CTLComponent', () => {
  let component: CTLComponent;
  let fixture: ComponentFixture<CTLComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CTLComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CTLComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
