import { Component, Inject } from '@angular/core';
import { Answer, Quiz, RoomService, Variant } from '../../../services/room.service';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { SocialUser } from '@abacritt/angularx-social-login';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatTableDataSource } from '@angular/material/table';

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
  series = [] as any

  showXAxis: boolean = true;
  showYAxis: boolean = true;
  gradient: boolean = true;
  showLegend: boolean = true;
  showXAxisLabel: boolean = true;
  xAxisLabel: string = 'ответы';
  showYAxisLabel: boolean = true;
  yAxisLabel: string = 'Количество ответов';
  legendTitle: string = 'Варианты ответа';

  tableData: {author:string, variantText:string}[] = []

  displayedColumns: string[] = ['author', 'variantText'];
  dataSource = new MatTableDataSource(this.tableData);

  mappedVariants: Map<number, string>

  constructor(public dialogRef: MatDialogRef<ShowQuizDialogComponent>,
              private roomService: RoomService,
              private snackBar: MatSnackBar,
              @Inject(MAT_DIALOG_DATA) public data: DialogData) {
    this.userAnswer = this.data.quiz.answers.find(v => v.author === this.data.user.email)

    const seriesCalc = new Map()

    this.mappedVariants = new Map(this.data.quiz.variants.map((v) => [v.id, v.value]))

    this.data.quiz.answers.forEach(v => {
      if (!seriesCalc.has(v.variant_id)) {
        seriesCalc.set(v.variant_id, 0)
      }
      seriesCalc.set(v.variant_id, seriesCalc.get(v.variant_id) + 1)
      this.tableData.push({author: v.author, variantText: this.mappedVariants.get(v.variant_id)!})
    })

    seriesCalc.forEach((v, k) => {
      this.series.push({name: this.mappedVariants.get(k)!, value: v})
    })


  }

  applyFilter(event: Event) {
    const filterValue = (event.target as HTMLInputElement).value;
    this.dataSource.filter = filterValue.trim().toLowerCase();
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
