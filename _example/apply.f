(apply (lambda (a) (times a 2))
  '(1.123 2.3332 3.14242 4 5.333)
)
// eat up comment
(func apply (f list)
  (while (not (empty list)) (f list) )
)
