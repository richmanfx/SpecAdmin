
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


// Подстановка имени сюиты в поле в модальном окне
// Для кнопки "Редактировать"
$('#editSuite').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);         // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');     // Извлечь информацию из "data-name" у кнопки
    var modal = $(this);         // Обновить модальное окно
    modal.find('#id_suite').val(recipientName)      // Только в input со старым значением
});

// Для кнопки "Удалить"
$('#delSuite').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);         // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');     // Извлечь информацию из "data-name" у кнопки
    var modal = $(this);         // Обновить модальное окно
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
    var button = $(event.relatedTarget);         // Кнопка, вызвавшая модальное окно
    var recipientName = button.data('name');     // Извлечь информацию из "data-name" у кнопки
    var recipientScriptsId = button.data('script');     // Извлечь "ScriptsId"

    var stepsScriptName = "kuku1";
    var scripsSuiteName = "hehe2";

    // Получить имя Сценария по его ScriptsId и имя Сюиты, которой принадлежит Сценарий
    $.ajax({
        type: 'POST',
        url: '/spec-admin/get-steps-options',
        data: 'ScriptsId=' + recipientScriptsId,
        success: function(answerFromServer){
//                alert( "Прибыли данные: " + answer_from_server );
//                (stepsScriptName = $(answer_from_server).find("stepsScriptName").text);
            var parsedData = JSON.parse(answerFromServer);
            alert( parsedData );
//                alert( "stepsScriptName=" + stepsScriptName );

        },
        error: function(){
            alert("Ошибка при ответе на AJAX POST запрос");
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
    var modal = $(this);         // Обновить модальное окно
    modal.find('#id_step').val(recipientName);
    modal.find('#id_steps_script').val(recipientScriptsId)      // В input с именем Сценария
});
