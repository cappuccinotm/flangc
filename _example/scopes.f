(setq x 15)
(cond (equal x 5) (prog () (
  (print 123)
  (return 45)
)) (prog () (
  (print 456)
  (return 12)
)))

(setq x 0)
(while (not (equal x 5)) (prog (x) (
  (setq x (plus x 1))
  (print x)
  (cond (equal x 2) (break))
))