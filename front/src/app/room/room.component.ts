import { Component, ElementRef, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { environment } from '../../environments/environment';
import { RoomService } from '../services/room.service';
import { filter } from 'rxjs';

@Component({
  selector: 'app-room',
  templateUrl: './room.component.html',
  styleUrls: ['./room.component.scss']
})
export class RoomComponent implements OnInit, OnDestroy {
  @ViewChild('localStream') localStream!: ElementRef<HTMLVideoElement>
  @ViewChild('remoteStream') remoteStream!: ElementRef<HTMLVideoElement>

  ws: WebSocket = {} as WebSocket // TODO: temporal
  roomID = 0
  userID = '' // TODO: temporal


  constructor(private route: ActivatedRoute, private _roomService: RoomService) {
    this.route.params.subscribe(params => {
      if (!!params['id']) {
        this.roomID = params['id']
      }
    })
  }

  ngOnInit(): void {
    this.userID = this._roomService.initPeer()
  }

  ngAfterViewInit(): void {
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
    // await this.impl2Service.establishMediaCall(this.userID)
    await this._roomService.establishMediaCall(String(this.roomID))
  }

  closeStream() {
    this._roomService.closeMediaCall()
  }


  ngOnDestroy(): void {
    this._roomService.destroyPeer();
  }

  conn() {
    this.ws = new WebSocket(`ws://localhost:${environment.port}/api/room/${this.roomID}`);

    this.ws.onopen = () => {
      this.ws.onmessage = (event) => {
        if (event.data == "ping") {
          this.ws.send("pong")
          return
        }
        console.log(event);
      }

      this.ws.send("Here's some text that the server is urgently awaiting!");
    };


  }

  sendMsg() {
    this.ws.send("random text");
  }
}
