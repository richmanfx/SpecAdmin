
// Скрипт для выпадающих списков
function anichange (objName) {
    if ( $(objName).css('display') === 'none' ) {
        $(objName).animate({height: 'show'}, 400);
    } else {
        $(objName).animate({height: 'hide'}, 200);
    }
}
function anichange(a) {
    $(a).slideToggle(400);
}


/// Подстановка имени сюиты в поле в модальном окне
// Для кнопки "Переименовать Сюиту"
$('#renameSuite').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);         // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');     // Извлечь информацию из "data-name" у кнопки
    var modal = $(this);                         // Обновить модальное окно
    modal.find('#id_old_suite').val(recipientName)      // Только в input со старым значением
});

// Для кнопки "Редактировать Сюиту"
$('#editSuite').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    var recipientName = button.data('name');
    var modal = $(this);
    modal.find('#id_suite').val(recipientName)
});

// Для кнопки "Удалить Сюиту"
$('#delSuite').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    var recipientName = button.data('name');
    var modal = $(this);
    modal.find('#id_suite').val(recipientName)
});

// Для кнопки "Добавить Сценарий"
$('#addScript').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    var recipientName = button.data('name');
    var modal = $(this);
    modal.find('#id_scripts_suite').val(recipientName)
});


// Подстановка имени Сценария и Сюиты в поля в модальном окне
// Для кнопки "Редактировать"
$('#editScript').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);         // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');     // Извлечь информацию из "data-name" у кнопки
    var recipientSuite = button.data('suite');     // Извлечь информацию из "data-suite" у кнопки
    var modal = $(this);         // Обновить модальное окно
    modal.find('#id_script').val(recipientName);      // В input с именем сценария
    modal.find('#id_scripts_suite').val(recipientSuite)      // В input с именем сюиты
});

// Для кнопки "Удалить"
$('#delScript').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);         // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');     // Извлечь информацию из "data-name" у кнопки
    var recipientSuite = button.data('suite');     // Извлечь информацию из "data-suite" у кнопки
    var modal = $(this);         // Обновить модальное окно
    modal.find('#id_script').val(recipientName);
    modal.find('#id_scripts_suite').val(recipientSuite)      // В input с именем сюиты
});

// Для кнопки "Добавить Шаг"
$('#addStep').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    var recipientName = button.data('name');
    var recipientSuite = button.data('suite');
    var modal = $(this);
    modal.find('#id_steps_script').val(recipientName);
    modal.find('#id_scripts_suite').val(recipientSuite)
});


// Подстановка имени Шага, его Сценария и Сюиты в поля в модальном окне
// Для кнопки "Редактировать"
$('#editStep').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);            // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');        // Извлечь информацию из "data-name" у кнопки
    var recipientScriptsId = button.data('script');     // Извлечь "ScriptsId"

    var stepsScriptName = "";
    var scripsSuiteName = "";

    // Получить имя Сценария по его ScriptsId и имя Сюиты, которой принадлежит Сценарий
    $.ajax({
        async: false,
        type: 'POST',
        url: '/spec-admin/get-steps-options',
        data: 'ScriptsId=' + recipientScriptsId,
        success: function(answerFromServer){
            // alert(answerFromServer.stepsScriptName + " и " + answerFromServer.scripsSuiteName);
            stepsScriptName = answerFromServer.stepsScriptName;
            scripsSuiteName = answerFromServer.scripsSuiteName;
        },
        error: function(){
            alert("Ошибка при ответе на AJAX POST запрос имени Сценария по его ScriptsId и имени Сюиты.");
        }

    });

    var modal = $(this);         // Обновить модальное окно
    modal.find('#id_step').val(recipientName);                      // В input с именем Шага
    modal.find('#id_steps_script').val(stepsScriptName);        // В input с именем Сценария
    modal.find('#id_scripts_suite').val(scripsSuiteName);       // В input с именем Сюиты

});

// Для кнопки "Удалить"
$('#delStep').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);         // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');     // Извлечь информацию из "data-name" у кнопки
    var recipientScriptsId = button.data('script');     // Извлечь "ScriptsId"

    var stepsScriptName = "";
    var scripsSuiteName = "";

    // Получить имя Сценария по его ScriptsId и имя Сюиты, которой принадлежит Сценарий
    $.ajax({
        async: false,
        type: 'POST',
        url: '/spec-admin/get-steps-options',
        data: 'ScriptsId=' + recipientScriptsId,
        success: function(answerFromServer){
            // alert(answerFromServer.stepsScriptName + " и " + answerFromServer.scripsSuiteName);
            stepsScriptName = answerFromServer.stepsScriptName;
            scripsSuiteName = answerFromServer.scripsSuiteName;
        },
        error: function(){
            alert("Ошибка при ответе на AJAX POST запрос удаления Шага");
        }

    });

    var modal = $(this);         // Обновить модальное окно
    modal.find('#id_step').val(recipientName);
    modal.find('#id_script').val(stepsScriptName);        // В input с именем Сценария
    modal.find('#id_suite').val(scripsSuiteName);       // В input с именем Сюиты
});


/// Печать на принтер
// Распечатать список шагов сценария
$('a#id_print_steps').on('click', function () {
    var scriptName = $(this).attr("data-name");    // Извлечь информацию из "data-name" у кнопки
    var suiteName = $(this).attr("data-suite");    // Извлечь информацию из "data-suite" у кнопки

    // alert(scriptName + " и " + suiteName);

    var scriptId = "";

    // Получить 'id' сценария по его имени и его сюите
    $.ajax({
        async: false,
        type: 'POST',
        url: '/spec-admin/print-scripts-steps',
        data: 'scriptName=' + scriptName + ';' + 'suiteName=' + suiteName,
        success: function(answerFromServer){
            console.log("success");
            alert(answerFromServer);

        },
        error: function(){
            alert("Ошибка при ответе на AJAX POST запрос для печати ПДФ");
        }
    });

    // // Получить список шагов сценария по его 'script_id'
    // $.ajax({
    //     async: false,
    //     type: 'POST',
    //     url: '/spec-admin/get-steps-from-script',
    //     data: 'scriptId=' + scriptId,
    //     success: function(answerFromServer){
    //         alert(answerFromServer)
    //     },
    //     error: function(){
    //         alert("Ошибка при ответе на AJAX POST запрос получения всех шагов сценария");
    //     }
    // });

    return false;
});














