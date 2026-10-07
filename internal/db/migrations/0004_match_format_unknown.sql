-- Allow matches whose series format is unknown (best_of = 0). Imported
-- tournament archives often only show the uploaded maps, which doesn't
-- prove whether a series was a Bo1 or a Bo3.
alter table matches drop constraint matches_best_of_check;
alter table matches add constraint matches_best_of_check check (best_of in (0, 1, 3, 5));
