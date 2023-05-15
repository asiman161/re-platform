import { Component, Inject } from '@angular/core';
import { Answer, Quiz, RoomService, Variant } from '../../../services/room.service';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { SocialUser } from '@abacritt/angularx-social-login';
import { MatSnackBar } from '@angular/material/snack-bar';

export interface DialogData {
  quiz: Quiz
  user: SocialUser
}

@Component({
  selector: 'app-show-quiz-dialog',
  templateUrl: './show-quiz-dialog.component.html',
  styleUrls: ['./show-quiz-dialog.component.scss']
})
export class ShowQuizDialogComponent {
  selectedVariant: Variant = {} as Variant
  userAnswer: Answer | undefined

  constructor(public dialogRef: MatDialogRef<ShowQuizDialogComponent>,
              private roomService: RoomService,
              private snackBar: MatSnackBar,
              @Inject(MAT_DIALOG_DATA) public data: DialogData) {
    this.userAnswer = this.data.quiz.answers.find(v => v.author === this.data.user.email)
  }

  selectVariant(variant: Variant) {
    this.selectedVariant = variant
  }

  closeQuiz() {
    this.roomService.closeQuiz(this.data.quiz.id).subscribe({
      next: () => this.dialogRef.close(),
      error: () => this.snackBar.open(`can't close quiz`, 'Close'),
    })
  }

  answerQuiz() {
    this.roomService.answerQuiz(this.data.quiz.id, this.selectedVariant).subscribe({
      next: () => this.dialogRef.close(),
      error: () => this.snackBar.open(`can't answer quiz`, 'Close'),
    })
  }
}
