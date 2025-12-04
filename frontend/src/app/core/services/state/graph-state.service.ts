import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { models } from '@wailsjs/go/models';
import { BaseStateService } from '@core/services/state/base-state.service';
import { AppConfigService } from '@core/config/app-config.service';

export type ViewMode = '2d' | '3d' | 'list';

export interface DisplayOptions {
  showIcons?: boolean;
  showRequired?: boolean;
  showOptional?: boolean;
  showInstalled?: boolean;
}

interface GraphState {
  graph: models.Graph | null;
  viewMode: ViewMode;
  displayOptions: DisplayOptions;
  selectedNodes: string[];
}

@Injectable({
  providedIn: 'root',
})
export class GraphStateService extends BaseStateService<GraphState> {
  private graph$ = new BehaviorSubject<models.Graph | null>(null);
  private viewMode$ = new BehaviorSubject<ViewMode>('list');
  private displayOptions$ = new BehaviorSubject<DisplayOptions>({
    showIcons: true,
    showRequired: true,
    showOptional: true,
    showInstalled: true,
  });
  private selectedNodes$ = new BehaviorSubject<string[]>([]);

  constructor(configService: AppConfigService) {
    super(
      'graph-state',
      {
        graph: null,
        viewMode: 'list',
        displayOptions: {
          showIcons: true,
          showRequired: true,
          showOptional: true,
          showInstalled: true,
        },
        selectedNodes: [],
      },
      configService
    );

    // Initialize from loaded state
    const state = this.getState();
    this.graph$.next(state.graph);
    this.viewMode$.next(state.viewMode);
    this.displayOptions$.next(state.displayOptions);
    this.selectedNodes$.next(state.selectedNodes);
  }

  setGraph(graph: models.Graph): void {
    this.graph$.next(graph);
    this.updateState({ graph });
  }

  setViewMode(mode: ViewMode): void {
    this.viewMode$.next(mode);
    this.updateState({ viewMode: mode });
  }

  updateDisplayOptions(options: Partial<DisplayOptions>): void {
    const newOptions = { ...this.displayOptions$.value, ...options };
    this.displayOptions$.next(newOptions);
    this.updateState({ displayOptions: newOptions });
  }

  selectNode(nodeId: string): void {
    const current = this.selectedNodes$.value;
    if (!current.includes(nodeId)) {
      const newSelection = [...current, nodeId];
      this.selectedNodes$.next(newSelection);
      this.updateState({ selectedNodes: newSelection });
    }
  }

  deselectNode(nodeId: string): void {
    const current = this.selectedNodes$.value;
    const newSelection = current.filter((id) => id !== nodeId);
    this.selectedNodes$.next(newSelection);
    this.updateState({ selectedNodes: newSelection });
  }

  clearSelection(): void {
    this.selectedNodes$.next([]);
    this.updateState({ selectedNodes: [] });
  }

  getGraph$(): Observable<models.Graph | null> {
    return this.graph$.asObservable();
  }

  getViewMode$(): Observable<ViewMode> {
    return this.viewMode$.asObservable();
  }

  getDisplayOptions$(): Observable<DisplayOptions> {
    return this.displayOptions$.asObservable();
  }

  getSelectedNodes$(): Observable<string[]> {
    return this.selectedNodes$.asObservable();
  }

  private updateState(partial: Partial<GraphState>): void {
    this.state$.next({ ...this.state$.value, ...partial });
  }
}

