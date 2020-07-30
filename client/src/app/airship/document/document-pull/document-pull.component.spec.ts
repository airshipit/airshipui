import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DocumentPullComponent } from './document-pull.component';

describe('DocumentPullComponent', () => {
  let component: DocumentPullComponent;
  let fixture: ComponentFixture<DocumentPullComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DocumentPullComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DocumentPullComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
