import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BareMetalComponent } from './bare-metal.component';

describe('BareMetalComponent', () => {
  let component: BareMetalComponent;
  let fixture: ComponentFixture<BareMetalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BareMetalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BareMetalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
