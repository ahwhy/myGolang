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
        return () =>{
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

function callback() {
    console.log('Done');
}

function testCallback() {
    console.log('before setTimeout() callback');
    setTimeout(callback, 1000); // 1秒钟后调用callback函数 
    console.log('after setTimeout() callback');
}

function testResultCallbackFunc(resolve, reject) {
    var timeOut = Math.random() * 2;
    console.log('set timeout to: ' + timeOut + ' seconds.');
    setTimeout(function () {
        if (timeOut < 1) {
            console.log('call resolve()...');
            resolve('200 OK');
        }
        else {
            console.log('call reject()...');
            reject('timeout in ' + timeOut + ' seconds.');
        }
    }, timeOut * 1000);
}

function testResultCallback() {
    success = (message) => {console.log(`success ${message}`)}
    failed = (error) => {console.log(`failed ${error}`)}
    testResultCallbackFunc(success, failed)
    Promise
}

function testWithPromise() {
    var p1 = new Promise(testResultCallbackFunc)
    p1.then((resp) => {
        console.log(resp)
    }).catch((err) => {
        console.log(err)
    })
}

// async/await 
async function testWithAsync() {
    var p1 = new Promise(testResultCallbackFunc)

    try {
        var resp = await p1
        console.log(resp)
    } catch (err) {
        console.log(err)
    }
}

function F() {
   this.name = '小明'
}

function testNewObj() {
    F()
    console.log('xx ', this.name)
}

function main() {
    testNewObj()
}

main()