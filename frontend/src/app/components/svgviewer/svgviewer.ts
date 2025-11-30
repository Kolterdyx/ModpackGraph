import { AfterViewInit, Component, ElementRef, Input, OnDestroy, ViewChild } from '@angular/core';
import SvgPanZoom from 'svg-pan-zoom';

@Component({
  selector: 'app-svgviewer',
  imports: [],
  templateUrl: './svgviewer.html',
  styleUrl: './svgviewer.scss',
})
export class SVGViewer implements AfterViewInit, OnDestroy {


  @ViewChild('svgContainer', {static: true}) svgContainer?: ElementRef<HTMLDivElement>;

  @Input() set svgData(data: string | undefined) {
    if (!data) return;
    const container = this.svgContainer?.nativeElement;
    if (!container) return;
    container.innerHTML = data;
    const svgEl = container.querySelector('svg');
    if (!svgEl) return;

    // Destroy old panzoom instance
    this.panZoom?.destroy();

    // Match viewport to container
    this.setSvgViewportToContainer(svgEl);

    // Init panzoom
    this.panZoom = SvgPanZoom(svgEl, {
      zoomEnabled: true,
      controlIconsEnabled: true,
      minZoom: 0.1,
      maxZoom: 10,
      fit: true,
      center: true,
    });

    // Re-fit at load
    setTimeout(() => {
      this.panZoom?.resize();
      this.panZoom?.fit();
      this.panZoom?.center();
    }, 0);
  };

  private panZoom?: SvgPanZoom.Instance;
  private resizeObserver?: ResizeObserver;

  constructor() {
    this.resizeObserver = new ResizeObserver(() => {
      this.updateViewportSize();
    });
    const container = this.svgContainer?.nativeElement;
    if (!container) return;
    this.resizeObserver.observe(container);
  }


  ngAfterViewInit() {
    // Observe container size changes
    this.resizeObserver = new ResizeObserver(() => {
      this.updateViewportSize();
    });
    const container = this.svgContainer?.nativeElement;
    if (!container) return;
    this.resizeObserver.observe(container);
  }

  ngOnDestroy() {
    this.resizeObserver?.disconnect();
    this.panZoom?.destroy();
  }

  private setSvgViewportToContainer(svgEl: SVGSVGElement) {
    const container = this.svgContainer?.nativeElement;
    if (!container) return;

    svgEl.setAttribute('width', `100%`);
    svgEl.setAttribute('height', `100%`);
    svgEl.setAttribute('preserveAspectRatio', `xMidYMid meet`);
    svgEl.setAttribute('preserveAspectRatio', 'xMidYMid meet');
    svgEl.style.width = '100%';
    svgEl.style.height = '100%';
  }

  private updateViewportSize() {
    if (!this.panZoom) return;

    const svgEl = this.svgContainer?.nativeElement?.querySelector('svg');
    if (!svgEl) return;

    // 1. Update viewport
    this.setSvgViewportToContainer(svgEl);

    // 2. Tell panzoom its container changed
    this.panZoom.resize();

    // 3. Re-fit the content to the resized container
    this.panZoom.fit();
    this.panZoom.center();
  }

}
