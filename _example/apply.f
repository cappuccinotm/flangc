(func apply (f list)
  (while (not (empty list)) (f list) )
)
// eat up comment
(apply (lambda (a) (a times 2)) '(1.123 2.3332 3.14242 4 5.333))