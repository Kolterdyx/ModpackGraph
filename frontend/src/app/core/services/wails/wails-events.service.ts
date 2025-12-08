import { Injectable } from '@angular/core';
import { Observable, fromEvent, Subject } from 'rxjs';
import { takeUntil, share } from 'rxjs/operators';
import { EventsOn, EventsOff, EventsEmit } from '@wailsjs/runtime/runtime';

export interface ProgressEvent {
  operation: string;
  message: string;
  progress: number;
}

@Injectable({
  providedIn: 'root',
})
export class WailsEventsService {
  private destroy$ = new Subject<void>();

  onProgress(): Observable<ProgressEvent> {
    return new Observable<ProgressEvent>((observer) => {
      const cleanup = EventsOn('progress', (data: ProgressEvent) => {
        observer.next(data);
      });

      return () => {
        cleanup();
      };
    }).pipe(takeUntil(this.destroy$), share());
  }

  onError(): Observable<any> {
    return new Observable<any>((observer) => {
      const cleanup = EventsOn('error', (data: any) => {
        observer.next(data);
      });

      return () => {
        cleanup();
      };
    }).pipe(takeUntil(this.destroy$), share());
  }

  on(eventName: string): Observable<any> {
    return new Observable<any>((observer) => {
      const cleanup = EventsOn(eventName, (data: any) => {
        observer.next(data);
      });

      return () => {
        cleanup();
      };
    }).pipe(takeUntil(this.destroy$), share());
  }

  emit(eventName: string, ...data: any[]): void {
    EventsEmit(eventName, ...data);
  }

  cleanup(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}

