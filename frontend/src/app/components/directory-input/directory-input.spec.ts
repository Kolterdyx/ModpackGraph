import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DirectoryInput } from './directory-input';

describe('DirectoryInput', () => {
  let component: DirectoryInput;
  let fixture: ComponentFixture<DirectoryInput>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DirectoryInput]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DirectoryInput);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
