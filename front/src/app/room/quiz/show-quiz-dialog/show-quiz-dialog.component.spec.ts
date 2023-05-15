import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShowQuizDialogComponent } from './show-quiz-dialog.component';

describe('ShowQuizDialogComponent', () => {
  let component: ShowQuizDialogComponent;
  let fixture: ComponentFixture<ShowQuizDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ShowQuizDialogComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ShowQuizDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
