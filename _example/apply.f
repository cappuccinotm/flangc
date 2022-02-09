(func apply (f list) (
  (while (not (empty list) (
    (f list)
  )))
))

(apply (lambda (a) (a times 2)) '(1 2 3 4 5)) // should return [2, 4, 6, 8, 10]