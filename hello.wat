(module
    (import "" "log" (func $log (param i32 i32)))
    (import "" "mem" (memory 1))
    (data (i32.const 0) "Hello, world!")
    (func (export "run")
        i32.const 0
        i32.const 13
        call $log
    )
)