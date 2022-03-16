(func listLengthHelper (list n)
    (cond (equal list '()) n
        (listLengthHelper (tail list) (plus n 1))
    )
)

(func listLength (list) (listLengthHelper list 0))

(listLength '(1 2 3 4 5 6))
