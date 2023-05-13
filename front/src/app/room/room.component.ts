import { Component, ElementRef, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Room, RoomService } from '../services/room.service';
import { WebsocketService } from '../services/websocket.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { AuthService } from '../services/auth.service';
import { SocialUser } from '@abacritt/angularx-social-login';

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
  room: Room = <Room>{}
  user: SocialUser = <SocialUser>{}


  constructor(private route: ActivatedRoute,
              private _roomService: RoomService,
              private _snackBar: MatSnackBar,
              private authService: AuthService) {
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


  ngOnInit(): void {
    this._roomService.getRoom(this.roomID).subscribe({
      next: room => {
        this.room = room
      },
      error: err => {
        this._snackBar.open(`can't get room ${err}`, "Close")
      }
    })
    // this.userID = this._roomService.initPeer(this.authService.user?.id || '')
  }

  ngAfterViewInit(): void {

  }

  ngOnDestroy(): void {
  }
}
