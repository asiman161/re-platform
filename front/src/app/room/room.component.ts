import { Component, ElementRef, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { RoomService } from '../services/room.service';
import { filter } from 'rxjs';
import { WebsocketService } from '../services/websocket.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { AuthService } from '../services/auth.service';

@Component({
  selector: 'app-room',
  templateUrl: './room.component.html',
  styleUrls: ['./room.component.scss']
})
export class RoomComponent implements OnInit, OnDestroy {
  @ViewChild('localStream') localStream!: ElementRef<HTMLVideoElement>
  @ViewChild('remoteStream') remoteStream!: ElementRef<HTMLVideoElement>

  ws: WebSocket = {} as WebSocket // TODO: temporal
  roomID = ''
  userID = '' // TODO: temporal


  constructor(private route: ActivatedRoute, private _roomService: RoomService, private _snackBar: MatSnackBar, private authService: AuthService, private websocketService: WebsocketService) {
    this.route.params.subscribe(params => {
      if (!!params['id']) {
        this.roomID = params['id']
      }
    })
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


  ngOnInit(): void {
    // this.userID = this._roomService.initPeer(this.authService.user?.id || '')
    this.userID = this._roomService.initPeer()
  }

  ngAfterViewInit(): void {

  }

  async waitStream() {
    //  this should be in ngAfterViewInit
    this._roomService.localStream$
      .pipe(filter(res => !!res))
      .subscribe(stream => {
        this.localStream.nativeElement.srcObject = stream
      })
    this._roomService.remoteStream$
      .pipe(filter(res => !!res))
      .subscribe((stream: any) => this.remoteStream.nativeElement.srcObject = stream)
  }

  async startStream() {
    await this._roomService.enableCallAnswer()
  }

  async joinStream() {
    console.log("pg: ", this.roomID)
    // await this._roomService.establishMediaCall(this.userID)
    await this._roomService.establishMediaCall(String(this.roomID))
  }

  closeStream() {
    this._roomService.closeMediaCall()
  }


  ngOnDestroy(): void {
    this._roomService.destroyPeer();
  }

  conn() {
    this.websocketService.conn(String(this.roomID))
  }

  sendMsg() {
    this.websocketService.sendMsg()
  }
}
