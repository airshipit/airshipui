import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {BaremetalComponent} from './baremetal.component';
import {MatButtonModule} from '@angular/material/button';
import {ToastrModule} from 'ngx-toastr';

describe('BaremetalComponent', () => {
  let component: BaremetalComponent;
  let fixture: ComponentFixture<BaremetalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        MatButtonModule,
        ToastrModule.forRoot()
      ],
      declarations: [
        BaremetalComponent
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BaremetalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
