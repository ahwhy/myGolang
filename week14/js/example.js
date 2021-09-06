function abs(x) {
    if (x >= 0) {
        return x;
    } else {
        return -x;
    }
}

function testMethod() {
    var person = {name: '小明', age: 23}
    person.greet = function() {
        console.log(`hello, my name is ${this.name}`)
    }

    person.greet()
}


function testMethod1() {
    fn = function() {
        console.log(this)
    }  
    fn()
}

function testMethod2() {
    var person = {name: '小明', age: 23}
    person.greetfn = function() {
        return function() {
            console.log(`hello, my name is ${this.name}`)
        }
    }

    person.greetfn()()
}




function ArrowFunction() {
    fn = x => x * x
    console.log(fn(10)) 
}

function testMethod3() {
    var person = {name: '小明', age: 23}
    person.greetfn = function() {
        return () => {
            console.log(`hello, my name is ${this.name}`)
        }
    }

    person.greetfn()()
}

function foo() {
    var x = 'Hello, ' + y;
    console.log(x);
    var y = 'Bob';
}

function main() {
    console.log(abs(-10))
    foo()
}

main()