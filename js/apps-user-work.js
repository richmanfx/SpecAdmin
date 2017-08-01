/**
 * Created by Александр Ящук (R5AM, Zoer) on 27.07.2017.
 */

/// Подстановка Логина и Полного имени пользователя в поля в модальном окне
// Для кнопки "Удалить"
$('#deleteUser').on('show.bs.modal', function (event) {

    // Кнопка, вызвавшая модальное окно
    var button = $(event.relatedTarget);

    // Извлечь информацию из "data-*" полей у кнопки
    var deletedLogin = button.data('login');
    var deletedFullName = button.data('name');

    // Обновить модальное окно
    var modal = $(this);

    // В input-ы вставить
    modal.find('#id_login').val(deletedLogin);
    modal.find('#id_full_name').val(deletedFullName);
});

// Для кнопки "Редактировать"
$('#editUser').on('show.bs.modal', function (event) {

    // Кнопка, вызвавшая модальное окно
    var button = $(event.relatedTarget);

    // Извлечь информацию из "data-*" полей у кнопки
    var editedLogin = button.data('login');
    var editedFullName = button.data('name');
    var editedPermissionCreate = button.data('perm-create');
    var editedPermissionEdit = button.data('perm-edit');
    var editedPermissionDelete = button.data('perm-delete');
    var editedPermissionConfig = button.data('perm-config');
    var editedPermissionUsers = button.data('perm-users');

    // Обновить модальное окно
    var modal = $(this);

    // В input-ы вставить
    modal.find('#id_login').val(editedLogin);
    modal.find('#id_full_name').val(editedFullName);

    // Чекбоксы выставить
    if (editedPermissionCreate === true) { modal.find('#id_create_permission').attr('checked', "checked") }
    if (editedPermissionEdit === true) { modal.find('#id_edit_permission').attr('checked', "checked") }
    if (editedPermissionDelete === true) { modal.find('#id_delete_permission').attr('checked', "checked") }
    if (editedPermissionConfig === true) { modal.find('#id_config_permission').attr('checked', "checked") }
    if (editedPermissionUsers === true) { modal.find('#id_users_permission').attr('checked', "checked") }
});


/// Вывод логина на label модальной формы изменения пароля
// Для кнопки "Изменить пароль"
$('#changePassword').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget);
    var Login = button.data('login');
    var modal = $(this);
    modal.find('#id_login').val(Login);     // В скрытый input
    document.getElementById('id_login_label').innerHTML = ' ' + Login;  // На label
});


/// Для исключения ошибки типа "An invalid form control with name='full_name' is not focusable." при
/// скрытии полей модальной формой.
jQuery(function ($) {
    $(document).on('nested:fieldRemoved', function (event) {
        $('[required]', event.field).removeAttr('required');
    });
});
