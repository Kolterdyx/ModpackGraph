import { ComponentFixture, TestBed } from '@angular/core/testing';

import { InteractiveTwoTab } from './interactive-two-tab';

describe('InteractiveTwoTab', () => {
  let component: InteractiveTwoTab;
  let fixture: ComponentFixture<InteractiveTwoTab>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [InteractiveTwoTab]
    })
    .compileComponents();

    fixture = TestBed.createComponent(InteractiveTwoTab);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
