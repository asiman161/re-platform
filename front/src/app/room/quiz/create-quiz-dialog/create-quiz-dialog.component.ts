import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { FormBuilder, Validators } from '@angular/forms';
import { Quiz, RoomService } from '../../../services/room.service';
import { MatSnackBar } from '@angular/material/snack-bar';


export interface DialogData {
  roomID: string
  author: string
}

@Component({
  selector: 'app-create-quiz-dialog',
  templateUrl: './create-quiz-dialog.component.html',
  styleUrls: ['./create-quiz-dialog.component.scss']
})
export class CreateQuizDialogComponent {
  quizForm = this.fb.group({
    name: this.fb.control('', [Validators.required, Validators.min(1)]),
    content: this.fb.control('', [Validators.required, Validators.min(1)]),
    variants: this.fb.array([
      this.fb.control('', [Validators.required, Validators.min(1)]),
      this.fb.control('', [Validators.required, Validators.min(1)]),
    ])
  });

  constructor(public dialogRef: MatDialogRef<CreateQuizDialogComponent>,
              private fb: FormBuilder,
              @Inject(MAT_DIALOG_DATA) public data: DialogData,
              private roomService: RoomService,
              private snackBar: MatSnackBar) {
  }

  addVariant() {
    const control = this.fb.control('', [Validators.required, Validators.min(1)])
    this.quizForm.controls.variants.push(control)
  }

  removeVariant(i: number) {
    if (this.quizForm.controls.variants.length > 1) {
      this.quizForm.controls.variants.removeAt(i)
    }
  }

  cancel() {
    this.dialogRef.close()
  }

  createQuiz() {
    if (this.quizForm.invalid) {
      this.snackBar.open(`Форма заполнена неверно`, "Close")
      return
    }

    const quiz: Partial<Quiz> = {
      room_id: this.data.roomID,
      author: this.data.author,
      name: this.quizForm.get("name")?.value!,
      content: this.quizForm.get("content")?.value!,
      variants: this.quizForm.controls.variants.controls.map((control, i) => {return {id: i, value: control.value!}}),
      is_open: true,
    }

    this.roomService.createQuiz(quiz).subscribe({
      next: () => this.dialogRef.close(),
      error: (err) => this.snackBar.open(`can't create quiz: ${err}`, 'Close'),
    })
  }
}
