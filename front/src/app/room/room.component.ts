import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';

@Component({
  selector: 'app-room',
  templateUrl: './room.component.html',
  styleUrls: ['./room.component.scss']
})
export class RoomComponent implements OnInit {
  ws: WebSocket = {} as WebSocket
  roomID = 0

  constructor(private route: ActivatedRoute, private http: HttpClient) {
    this.route.params.subscribe(params => {
      if (!!params['id']) {
        this.roomID = params['id']
      }
    })
  }

  ngOnInit(): void {
  }

  ping() {
    this.http.get(`http://localhost:${environment.port}/api/ping`, { responseType: 'text' }).subscribe({
      next: data => {
        console.log(data)
      },
      error: err => {
        console.error(err)
      }
    })
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
