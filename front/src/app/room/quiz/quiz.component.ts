import { Component, Input } from '@angular/core';
import { Quiz } from '../../services/room.service';
import { MatDialog } from '@angular/material/dialog';
import { DialogData, ShowQuizDialogComponent } from './show-quiz-dialog/show-quiz-dialog.component';
import { SocialUser } from '@abacritt/angularx-social-login';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-quiz',
  templateUrl: './quiz.component.html',
  styleUrls: ['./quiz.component.scss']
})
export class QuizComponent {
  @Input() quizzes!: Quiz[]
  user: SocialUser

  constructor(public dialog: MatDialog, private authService: AuthService) {
    this.user = this.authService.user!
  }

  openDetail(quiz: Quiz) {
    const dialogData: DialogData = { quiz: quiz, user: this.user }

    this.dialog.open(ShowQuizDialogComponent, {
      width: '400px',
      data: dialogData
    })
  }
}
