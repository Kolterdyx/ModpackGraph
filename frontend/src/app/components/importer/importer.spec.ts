import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Importer } from './importer';

describe('Importer', () => {
  let component: Importer;
  let fixture: ComponentFixture<Importer>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Importer]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Importer);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
