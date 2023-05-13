import { Component, Input, OnInit } from '@angular/core';
import { WebsocketService } from '../../services/websocket.service';
import { ChatMessage, ChatService } from '../../services/chat.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { SocialUser } from '@abacritt/angularx-social-login';
import { AuthService } from '../../services/auth.service';
import { FormBuilder, FormControl, Validators } from '@angular/forms';

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit {
  @Input() roomID = ''
  @Input() user!: SocialUser

  chatInput = new FormControl('', [Validators.required, Validators.min(1)])
  messages: ChatMessage[] = []
  chatForm = this._formBuilder.group({
    msg: this.chatInput
  });

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

  conn() {
    this.websocketService.conn(String(this.roomID))
  }

  sendMsg() {
    if (this.chatForm.valid) {
      this.websocketService.sendMsg(this.chatForm.get("msg")?.value!)
      this.chatInput.setValue("")
    }
  }
}
