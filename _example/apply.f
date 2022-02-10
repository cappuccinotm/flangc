(func apply (f list) (
  (while (not (empty list) (
    (f list)
  )))
))

(apply (lambda (a) (a times 2)) '(1 2 3 4 5))