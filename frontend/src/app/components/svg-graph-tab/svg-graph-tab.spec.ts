import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SvgGraphTab } from './svg-graph-tab';

describe('SvgGraphTab', () => {
  let component: SvgGraphTab;
  let fixture: ComponentFixture<SvgGraphTab>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SvgGraphTab]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SvgGraphTab);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
