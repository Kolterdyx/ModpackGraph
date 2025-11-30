import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SVGViewer } from './svgviewer';

describe('SVGViewer', () => {
  let component: SVGViewer;
  let fixture: ComponentFixture<SVGViewer>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SVGViewer]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SVGViewer);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
