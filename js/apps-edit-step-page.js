
// Удаление скриншота при нажатии ссылки с изображением красного крестика
$('#id_del_screen_shot').on('click', function () {
    var stepsId = $(this).attr("data-stepId");
    // alert(stepsId);

    // Удалить скриншот у заданного по Id Шага
    $.ajax({
        async: false,
        type: 'POST',
        url: '/spec-admin/del-screen-shot',
        data: 'StepsId=' + stepsId,
        success: function(answerFromServer){
            var deleteStatus = answerFromServer.deleteStatus;

            // Получить статус "OK" и обновить на странице DIV с изображением скриншота
            if (deleteStatus === "OK") {
                location.reload(true);
            }
        },
        error: function(){
        alert("Ошибка при ответе на AJAX POST запрос на удаление из Шага скриншота.");
    }
    });

    return false;
});
