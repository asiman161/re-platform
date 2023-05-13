import { Component, OnInit } from '@angular/core';
import { AppService } from '../services/app.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { RoomService } from '../services/room.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  constructor(private router: Router, private _appService: AppService, private roomService: RoomService, private _snackBar: MatSnackBar) {
  }

  ping() {
    this._appService.ping().subscribe({
      next: value => {
        this._snackBar.open(value, "Close")
      },
      error: () => {
        this._snackBar.open("Failed to ping", "Close")
      }
    })
  }

  ngOnInit(): void {
  }

  createRoom(): void {
    this.roomService.createRoom("new custom room").subscribe({
      next: room => {
        this.router.navigateByUrl(`/room/${room.id}`)
      },
      error: () => {
        this._snackBar.open("Failed to create room", "Close")
      }
    })
  }

  getRooms(): void {
    this.roomService.getRooms().subscribe({
      next: rooms => {
        console.log(rooms)
      },
      error: () => {
        this._snackBar.open("Failed to create room", "Close")
      }
    })
  }
}
