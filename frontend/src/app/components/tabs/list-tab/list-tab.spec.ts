import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ListTab } from './list-tab';

describe('ListTab', () => {
  let component: ListTab;
  let fixture: ComponentFixture<ListTab>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ListTab]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ListTab);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
