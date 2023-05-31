import { Component, ElementRef, HostListener, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Quiz, Room, RoomService, UserActivity } from '../services/room.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { AuthService } from '../services/auth.service';
import { SocialUser } from '@abacritt/angularx-social-login';
import { MatDialog } from '@angular/material/dialog';
import { CreateQuizDialogComponent, DialogData } from './quiz/create-quiz-dialog/create-quiz-dialog.component';
import { ChatMessage, ChatService } from '../services/chat.service';
import { WebSocketSubject } from 'rxjs/webSocket';
import { WebsocketService } from '../services/websocket.service';

@Component({
  selector: 'app-room',
  templateUrl: './room.component.html',
  styleUrls: ['./room.component.scss']
})
export class RoomComponent implements OnInit, OnDestroy {
  @ViewChild('localStream') localStream!: ElementRef<HTMLVideoElement>
  @ViewChild('remoteStream') remoteStream!: ElementRef<HTMLVideoElement>

  @HostListener('window:visibilitychange') lostFocus() {
    this.changeVisibility()
  }

  roomID = ''
  room: Room = <Room>{}
  user: SocialUser = <SocialUser>{}
  users: UserActivity[] = []
  messages: ChatMessage[] = []
  quizzes: Quiz[] = []
  // @ts-ignore
  ws: WebSocketSubject<any>


  constructor(private route: ActivatedRoute,
              private router: Router,
              private _roomService: RoomService,
              private _snackBar: MatSnackBar,
              private chatService: ChatService,
              private authService: AuthService,
              private websocketService: WebsocketService,
              public dialog: MatDialog) {
    this.route.params.subscribe(params => {
      if (!!params['id']) {
        this.roomID = params['id']
      }
    })

    this.user = this.authService.user!
  }


  closeRoom(): void {
    this._roomService.closeRoom(this.roomID).subscribe({
      next: () => {
        this._snackBar.open(`room with id ${this.roomID} closed`, "Close")
      },
      error: () => {
        this._snackBar.open(`Failed to close room ${this.roomID}`, "Close")
      }
    })
  }


  async ngOnInit() {
    this.getRoom()
    this.getMessages()
    this.getMessages()
    this.getQuizzes()
    this.changeVisibility()
    await this.conn()
  }

  getRoom() {
    this._roomService.getRoom(this.roomID).subscribe({
      next: room => {
        this.room = room
      },
      error: err => {
        this._snackBar.open(`can't get room ${err}`, "Close")
      }
    })
  }

  getQuizzes() {
    this._roomService.getQuizzes(this.roomID).subscribe({
      next: quizzes => this.quizzes = quizzes,
      error: err => this._snackBar.open(`Failed to get quizzes from room: ${this.roomID}. ${err}`, "Close")
    })
  }

  getMessages() {
    this.chatService.getMessages(this.roomID).subscribe({
      next: messages => this.messages = messages,
      error: err => this._snackBar.open(`Failed to get messages from room: ${this.roomID}. ${err}`, "Close")
    })
  }

  changeVisibility() {
    this._roomService.changeVisibility(this.roomID,{connected: true, active: !document.hidden}).subscribe({
      error: err => this._snackBar.open(`Failed to update activity status: ${this.roomID}. ${err}`, "Close")
    })
  }

  async conn() {
    this.ws = this.websocketService.conn(this.roomID, this.user.email)

    this.ws.subscribe({
      next: async v => {
        switch (v.type) {
          case 'message':
            // hack to force angular update
            this.messages.push(JSON.parse(v.data))
            this.messages = ([] as ChatMessage[]).concat(this.messages)
            break
          case 'quiz':
          case 'new_answer':
          case 'close_quiz':
            this.getQuizzes()
            break
          case 'close_room':
            await this.router.navigateByUrl('/')
            break
          case 'change_visibility':
            this.getRoomUsers()
            break
          default:
            this._snackBar.open(`unknown websocket message type: ${v.type}`, 'close')
            break
        }
      }
    })
  }

  ngAfterViewInit(): void {

  }

  getRoomUsers() {
    this._roomService.getRoomUsers(this.roomID).subscribe({
      next: value => this.users = value,
      error: err => this._snackBar.open(`can't get room users ${err}`, 'close')
    })
  }

  createQuizDialog() {
    const dialogData: DialogData = {
      roomID: this.roomID,
      author: this.user.email,
    }

    const dialogRef = this.dialog.open(CreateQuizDialogComponent, {
      width: "400px",
      data: dialogData
    })

    dialogRef.afterClosed().subscribe(() => {
    });
  }

  ngOnDestroy(): void {
    this.ws.complete()
  }
}
