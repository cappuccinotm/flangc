(func apply (f list)
  (cond (empty list) (
    return
  )(
    cons (f (head list)) (apply f (tail list))
  ))
)
// eat up comment
(print (apply (lambda (a) (times a 2))
         '(1.123 2.3332 3.14242 4 5.333)
       ))