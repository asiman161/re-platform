import { Injectable } from '@angular/core';
import { BehaviorSubject, Subject } from 'rxjs';
import { MediaConnection, Peer, PeerJSOption } from 'peerjs';
import { v4 as uuidv4 } from 'uuid';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class RoomService {

  private peer?: Peer;
  private mediaCall: MediaConnection = <MediaConnection>{};

  private localStreamBs: BehaviorSubject<MediaStream> = new BehaviorSubject(new MediaStream());
  public localStream$ = this.localStreamBs.asObservable();
  private remoteStreamBs: any = new BehaviorSubject(null);
  public remoteStream$ = this.remoteStreamBs.asObservable();

  private isCallStartedBs = new Subject<boolean>();
  public isCallStarted$ = this.isCallStartedBs.asObservable();

  public initPeer(): string {
    if (!this.peer || this.peer.disconnected) {
      const peerJsOptions: PeerJSOption = {
        debug: 3,
        host: environment.host,
        port: 3000,
        path: 'peerjs/myapp'
      };
      try {
        let id = uuidv4();
        this.peer = new Peer(id, peerJsOptions);
        return id;
      } catch (error) {
        console.error(error);
      }
    }

    return ''
  }

  public async establishMediaCall(remotePeerId: string) {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });

      const connection = this.peer!.connect(remotePeerId);
      connection.on('error', err => {
        console.error(err);
      });

      this.mediaCall = this.peer!.call(remotePeerId, stream);
      if (!this.mediaCall) {
        let errorMessage = 'Unable to connect to remote peer';
        throw new Error(errorMessage);
      }
      this.localStreamBs.next(stream);
      this.isCallStartedBs.next(true);

      this.mediaCall.on('stream',
        (remoteStream) => {
          this.remoteStreamBs.next(remoteStream);
        });
      this.mediaCall.on('error', err => {
        // this.snackBar.open(err, 'Close');
        console.error(err);
        this.isCallStartedBs.next(false);
      });
      this.mediaCall.on('close', () => this.onCallClose());
    } catch (ex) {
      console.error(ex);
      // this.snackBar.open(ex, 'Close');
      this.isCallStartedBs.next(false);
    }
  }

  public async enableCallAnswer() {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
      this.localStreamBs.next(stream);
      this.peer!.on('call', async (call) => {

        this.mediaCall = call;
        this.isCallStartedBs.next(true);

        this.mediaCall.answer(stream);
        this.mediaCall.on('stream', (remoteStream) => {
          this.remoteStreamBs.next(remoteStream);
        });
        this.mediaCall.on('error', err => {
          // this.snackBar.open(err, 'Close');
          this.isCallStartedBs.next(false);
          console.error(err);
        });
        this.mediaCall.on('close', () => this.onCallClose());
      });
    } catch (ex) {
      console.error(ex);
      // this.snackBar.open(ex, 'Close');
      this.isCallStartedBs.next(false);
    }
  }

  public closeMediaCall() {
    this.mediaCall?.close();
    if (!this.mediaCall) {
      this.onCallClose()
    }
    this.isCallStartedBs.next(false);
  }

  private onCallClose() {
    this.remoteStreamBs?.value.getTracks().forEach((track:any) => {
      track.stop();
    });
    this.localStreamBs?.value.getTracks().forEach(track => {
      track.stop();
    });
    // this.snackBar.open('Call Ended', 'Close');
  }

  public destroyPeer() {
    this.mediaCall?.close();
    this.peer?.disconnect();
    this.peer?.destroy();
  }
}
