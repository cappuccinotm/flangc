(setq x 15)
(func test (y) (prog (x) (
  (print (plus x 5))
  (print (minus y 5))
  (return (plus x y))
  (print (minus x y))
)))

// 20
// 5
// 25
(print (test 10))

(setq x 0)
(while (not (equal x 5)) (prog (x) (
  (setq x (plus x 1))
  (print x)
  (cond (equal x 2) (break))
)))