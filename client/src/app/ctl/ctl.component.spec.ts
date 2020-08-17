import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {CtlComponent} from './ctl.component';
import {RouterTestingModule} from '@angular/router/testing';

describe('CtlComponent', () => {
  let component: CtlComponent;
  let fixture: ComponentFixture<CtlComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        RouterTestingModule
      ],
      declarations: [CtlComponent]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CtlComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
