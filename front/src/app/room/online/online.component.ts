import { Component, Input } from '@angular/core';

export interface ActiveUser {
  email: string
  active: boolean
  connected: boolean
  created_at: string
}

@Component({
  selector: 'app-online',
  templateUrl: './online.component.html',
  styleUrls: ['./online.component.scss']
})
export class OnlineComponent {
  @Input() roomID = ''

  @Input() users: ActiveUser[] = []

  filterConnected(users: ActiveUser[]): ActiveUser[] {
    return users.filter(u => u.connected)
  }
}
