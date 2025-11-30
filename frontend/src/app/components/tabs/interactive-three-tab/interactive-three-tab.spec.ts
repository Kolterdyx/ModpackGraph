import { ComponentFixture, TestBed } from '@angular/core/testing';

import { InteractiveThreeTab } from './interactive-three-tab';

describe('InteractiveThreeTab', () => {
  let component: InteractiveThreeTab;
  let fixture: ComponentFixture<InteractiveThreeTab>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [InteractiveThreeTab]
    })
    .compileComponents();

    fixture = TestBed.createComponent(InteractiveThreeTab);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
