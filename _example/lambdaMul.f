(func lambdaMul ()
    (lambda (x)
        (lambda (y) (times x y))
    )
)

((lambdaMul 4) 8)
