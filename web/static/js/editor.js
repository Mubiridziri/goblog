(function() {
    var simplemde = new SimpleMDE({
        element: document.getElementById("editor"),
        placeholder: "Привет! Моя статья о...",
    });
    const tags = UseBootstrapTag(document.getElementById('tags'));

    // const form = document.getElementById("article-form")

//     form.onsubmit = function(event) {
//         event.preventDefault()
//         const values =new URLSearchParams(new FormData(this))
//         console.log(values)
//
//         fetch("/articles/new",
//             {   method: 'POST',
//                 mode : 'same-origin',
//                 credentials: 'same-origin' ,
//                 body : values
//             })
//             .then(function(response) {
//
//                 const toastTrigger = document.getElementById('liveToastBtn')
//                 const toastLiveExample = document.getElementById('liveToast')
//
//                 if (response.ok) {
//                     if (toastTrigger) {
//                         const toastBootstrap = bootstrap.Toast.getOrCreateInstance(toastLiveExample)
//                         toastTrigger.addEventListener('click', () => {
//                             toastBootstrap.show()
//                         })
//                     }
//                 }
//             });
//     };
//
//
//
})();