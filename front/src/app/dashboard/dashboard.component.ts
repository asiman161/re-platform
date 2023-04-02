import { Component, OnInit } from '@angular/core';
import { AppService } from '../services/app.service';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  constructor(private _appService: AppService, private _snackBar: MatSnackBar) {}

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

}
