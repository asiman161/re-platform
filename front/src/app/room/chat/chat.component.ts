import { Component, ElementRef, Input, OnDestroy, OnInit, QueryList, ViewChild, ViewChildren } from '@angular/core';
import { WebsocketService } from '../../services/websocket.service';
import { ChatMessage, ChatService } from '../../services/chat.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { SocialUser } from '@abacritt/angularx-social-login';
import { AuthService } from '../../services/auth.service';
import { FormBuilder, FormControl, Validators } from '@angular/forms';
import { WebSocketSubject } from 'rxjs/webSocket';

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit, OnDestroy {
  @Input() roomID = ''
  @Input() user!: SocialUser

  @ViewChildren('messagesElem') messagesElem!: QueryList<any>;
  @ViewChild('content') content!: ElementRef;

  chatInput = new FormControl('', [Validators.required, Validators.min(1)])
  messages: ChatMessage[] = []
  chatForm = this._formBuilder.group({
    msg: this.chatInput
  });
  // @ts-ignore
  ws: WebSocketSubject<any>

  constructor(private authService: AuthService,
              private snackBar: MatSnackBar,
              private websocketService: WebsocketService,
              private chatService: ChatService,
              private _formBuilder: FormBuilder) {
  }

  ngOnInit(): void {
    this.conn()
    this.chatService.getMessages(this.roomID).subscribe({
      next: messages => {
        this.messages = messages
      },
      error: err => {
        this.snackBar.open(`Failed to get messages from room: ${this.roomID}. ${err}`, "Close")
      }
    })
  }

  ngAfterViewInit() {
    this.scrollToBottom();
    this.messagesElem.changes.subscribe(this.scrollToBottom);
  }

  scrollToBottom = () => {
    try {
      this.content.nativeElement.scrollTop = this.content.nativeElement.scrollHeight;
    } catch (err) {}
  }

  conn() {
    this.ws = this.websocketService.conn(String(this.roomID))

    this.ws.subscribe({
      next: v => {
        this.messages.push(v)
      }
    })
  }

  sendMsg() {
    if (this.chatForm.valid) {
      this.websocketService.sendMsg(this.chatForm.get("msg")?.value!)
      this.chatInput.setValue("")
    }
  }

  ngOnDestroy(): void {
    this.ws.complete()
  }



}
