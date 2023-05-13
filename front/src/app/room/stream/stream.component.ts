import { Component, ElementRef, ViewChild } from '@angular/core';
import { StreamService } from '../../services/stream.service';
import { MatDialog } from '@angular/material/dialog';
import { filter, Observable, of, switchMap } from 'rxjs';
import { CallInfoDialogComponents, DialogData } from './dialog/callinfo-dialog.component';

@Component({
  selector: 'app-stream',
  templateUrl: './stream.component.html',
  styleUrls: ['./stream.component.scss']
})
export class StreamComponent {
  public isCallStarted$: Observable<boolean>;
  private peerId: string;
  showLocalVideo = false;

  @ViewChild('localVideo') localVideo!: ElementRef<HTMLVideoElement>;
  @ViewChild('remoteVideo') remoteVideo!: ElementRef<HTMLVideoElement>;

  constructor(public dialog: MatDialog, private callService: StreamService) {
    this.isCallStarted$ = this.callService.isCallStarted$;
    this.peerId = this.callService.initPeer();
  }

  ngOnInit(): void {
    this.callService.localStream$
      .pipe(filter(res => !!res))
      .subscribe(stream => this.localVideo.nativeElement.srcObject = stream)
    this.callService.remoteStream$
      .pipe(filter(res => !!res))
      .subscribe(stream => this.remoteVideo.nativeElement.srcObject = stream)
  }

  ngOnDestroy(): void {
    this.callService.destroyPeer();
  }

  public showModal(joinCall: boolean): void {
    // @ts-ignore
    let dialogData: DialogData = joinCall ? ({ peerId: null, joinCall: true }) : ({
      peerId: this.peerId,
      joinCall: false
    });
    const dialogRef = this.dialog.open(CallInfoDialogComponents, {
      width: '250px',
      data: dialogData
    });

    dialogRef.afterClosed()
      .pipe(
        switchMap(peerId =>
          joinCall ? of(this.callService.establishMediaCall(peerId)) : of(this.callService.enableCallAnswer())
        ),
      )
      .subscribe(_ => {
        this.showLocalVideo = true
      });
  }

  public endCall() {
    this.callService.closeMediaCall();
    this.showLocalVideo = false
  }
}
